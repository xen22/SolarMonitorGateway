#!/usr/bin/python

import time
import os
import datetime
import logging
import logging.handlers
import requests
import RPi.GPIO as GPIO

import gdrive

######################################################################################
# Constants
######################################################################################
class Constants(object):
    CURL_CMD = "/usr/bin/curl"
    GDRIVE_CMD = "/usr/local/bin/gdrive"
    GPIO_CMD = "/usr/bin/gpio"
    GDRIVE_CACHE_DIR = "/home/sol/.gdrive_cache"
    CAMERA_IP = "192.168.100.108"
    CAMERA_USERNAME = "admin"
    CAMERA_PWD = "**"
    CAMERA_INITIAL_WAIT_SECS = 0
    CAMERA_CONNECT_TIMEOUT_SECS = 3 * 60
    NIGHTTIME_WAIT_SECS = 10
    DAYTIME_WAIT_SECS = 3
    SEC_TO_WAIT_FOR_UPLOAD = 0
    SNAPSHOT_URL = "http://{ip}/cgi-bin/snapshot.cgi".format(ip=CAMERA_IP)
    SNAPSHOT_DIR = "/home/sol/camera_snapshots"
    SWITCH_GPIO_PIN = 26 
    CAMERA_GPIO_PIN = 16
    LOG_FILENAME = "/home/sol/take_snapshot.log"


######################################################################################
# Logging setup
######################################################################################
logging.basicConfig(level=logging.DEBUG,
                    format='%(asctime)s %(levelname)-8s %(message)s',
                    datefmt='%Y-%m-%d %H:%M:%S',
                    filename=Constants.LOG_FILENAME)

# Set up a specific logger with our desired output level
my_logger = logging.getLogger('MyLogger')
my_logger.setLevel(logging.DEBUG)

# Add the log message handler to the logger
handler = logging.handlers.RotatingFileHandler(
              Constants.LOG_FILENAME, maxBytes=100000, backupCount=5)
my_logger.addHandler(handler)

######################################################################################
# Helper functions
######################################################################################
def logGpio():
    #logging.debug(os.popen("{gpio} readall | grep GPIO.2[12]".format(gpio = Constants.GPIO_CMD)).read())
    return

def cleanup():
    logging.info("  Turning camera off.")
    GPIO.output(Constants.CAMERA_GPIO_PIN, GPIO.LOW)
    logging.info("  Turning switch off.")
    GPIO.output(Constants.SWITCH_GPIO_PIN, GPIO.LOW)
    logGpio()
    return

######################################################################################
# Main script
######################################################################################

try:
    logging.info("================================= BEGIN ===============================================")

    startTime = datetime.datetime.now()

    GPIO.setmode(GPIO.BCM)
    pinList = [Constants.SWITCH_GPIO_PIN, Constants.CAMERA_GPIO_PIN]
    logGpio()

    for i in pinList: 
        GPIO.setup(i, GPIO.OUT) 
        GPIO.output(i, GPIO.LOW)
    logGpio()

    if not os.path.exists(Constants.SNAPSHOT_DIR):
        os.makedirs(Constants.SNAPSHOT_DIR)

    logging.info("(1) Turning on switch.")
    GPIO.output(Constants.SWITCH_GPIO_PIN, GPIO.HIGH)
    logGpio()

    logging.info("(2) Turning on camera.")
    GPIO.output(Constants.CAMERA_GPIO_PIN, GPIO.HIGH)
    logGpio()

    logging.info("(3) Waiting for camera to become available.")
    time.sleep(Constants.CAMERA_INITIAL_WAIT_SECS)

    cameraRebooting = True
    while cameraRebooting:
        try:
            logging.debug("About to connect to " + Constants.SNAPSHOT_URL)
            requests.get(Constants.SNAPSHOT_URL, auth=(Constants.CAMERA_USERNAME, Constants.CAMERA_PWD))
            cameraRebooting = False
        except requests.exceptions.ConnectionError:
            logging.debug("conn error")
            if (datetime.datetime.now() - startTime).total_seconds() > Constants.CAMERA_CONNECT_TIMEOUT_SECS:
                logging.info("ERROR: timed out while waiting for the camera to become available ({}). Turning off camera and switch and exiting. ".format(Constants.CAMERA_CONNECT_TIMEOUT_SECS))
                cleanup()
                logging.fatal("Exiting unexpectedly.")
                exit(1)
            time.sleep(0.5) 
        except requests.exceptions.Timeout:
            logging.debug("timeout")

    now = datetime.datetime.now()

    # we need to wait longer if it's night time
    if now.hour >= 0 and now.hour <= 5:
        logging.info("(3.1): Waiting for the camera to adjust to the IR mode.")
        time.sleep(Constants.NIGHTTIME_WAIT_SECS)
    else:
        logging.info("(3.1): Waiting for the camera to adjust (normal daytime mode).")
        time.sleep(Constants.DAYTIME_WAIT_SECS)

    logging.info("(4) Downloading image from camera.")
    image_filename="{name}.jpeg".format(name=now.strftime("%Y.%m.%d_%H:%M:%S"))
    os.system("{curl} --user {user}:{pwd} --max-time 10 -o {dir}/{file} {url}".format(curl=Constants.CURL_CMD, 
    user=Constants.CAMERA_USERNAME, pwd=Constants.CAMERA_PWD, dir=Constants.SNAPSHOT_DIR, file=image_filename, url=Constants.SNAPSHOT_URL))

    logging.info("(5) Turning camera off.")
    GPIO.output(Constants.CAMERA_GPIO_PIN, GPIO.LOW)
    logGpio()

    logging.info("(6) Uploading image to Google Drive.")
    currentYear = now.strftime("%Y")
    currentMonth = now.strftime("%b")
    currentDay = now.strftime("%d")
    folder_id = ""
    yearFolderId = ""

    gdrive = gdrive.GDrive()
    gdrive.pushFile(os.path.join(Constants.SNAPSHOT_DIR, image_filename), "Cameras/{year}/{month}/{day}".format(year = currentYear, month = currentMonth, day = currentDay)) 

    time.sleep(Constants.SEC_TO_WAIT_FOR_UPLOAD)

    logging.info("(7) Turning switch off.")
    GPIO.output(Constants.SWITCH_GPIO_PIN, GPIO.LOW)
    logGpio()

    logging.info("Total duration: {}".format(datetime.datetime.now() - startTime))
    logging.info("================================= END =================================================")

except Exception, e:
    logging.error("Exception caught: {}".format(str(e)))
    cleanup()
    logging.fatal("Exiting unexpectedly.")
    exit(1)
