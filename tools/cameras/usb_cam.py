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
    USB_CAMERA_DEV = "/dev/video0"
    USB_CAMERA_RESOLUTION = "1280x960"
    USB_CAMERA_INITIAL_WAIT_SECS = 1
    USB_CAMERA_FRAMES_TO_SKIP = 3
    USB_CAMERA_NAME = "Solar shed"
    SNAPSHOT_DIR = "/home/sol/camera_snapshots"
    FSWEBCAM_CMD = "/usr/bin/fswebcam"
    LIGHT_GPIO_PIN = 6
    JPEG_QUALITY = 75
    SEC_TO_WAIT_FOR_UPLOAD = 1
    LOG_FILENAME = "/home/sol/take_snapshot.log"
    GPIO_CMD = "/usr/bin/gpio"


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

def initGPIO():
    GPIO.setmode(GPIO.BCM)
    pinList = [Constants.LIGHT_GPIO_PIN]
    logGpio()

    for i in pinList: 
        GPIO.setup(i, GPIO.OUT) 
        GPIO.output(i, GPIO.LOW)
    logGpio()

def logGpio():
    logging.debug(os.popen("{gpio} readall | grep GPIO.22".format(gpio=Constants.GPIO_CMD)).read())
    return

def lightOn(on):
    if on:
        logging.info("Turning on light.")
        GPIO.output(Constants.LIGHT_GPIO_PIN, GPIO.HIGH)
    else:
        logging.info("Turning off light.")
        GPIO.output(Constants.LIGHT_GPIO_PIN, GPIO.LOW)
    logGpio()

def cleanup():
    logging.info("Removing uvcvideo kernel module to disable camera.")
    os.system("sudo /sbin/rmmod uvcvideo")
    os.system("/sbin/lsmod | grep uvcvideo")
    lightOn(False)
    return

######################################################################################
# Main script
######################################################################################

try:
    logging.info("================================= BEGIN ===============================================")

    startTime = datetime.datetime.now()

    initGPIO()

    if not os.path.exists(Constants.SNAPSHOT_DIR):
        os.makedirs(Constants.SNAPSHOT_DIR)

    lightOn(True)

    logging.info("Turning on USB camera.")
    os.system("sudo /sbin/modprobe uvcvideo")
    os.system("/sbin/lsmod | grep uvcvideo")

    logging.info("Waiting for camera to become available.")
    time.sleep(Constants.USB_CAMERA_INITIAL_WAIT_SECS)

    now = datetime.datetime.now()
    logging.info("Downloading image from camera.")
    image_filename = "{name}_C2.jpeg".format(name=now.strftime("%Y.%m.%d_%H:%M:%S"))
    os.system("{fswebcam} -d {dev} -r {res} -S {skip} --jpeg {jpeg_qual} -D 1 --font \"sans:22\" --title \"{cam_name}\" --timestamp \"{time}\" {dir}/{file}".format(
        fswebcam=Constants.FSWEBCAM_CMD,
        dev=Constants.USB_CAMERA_DEV,
        res=Constants.USB_CAMERA_RESOLUTION,
        skip=Constants.USB_CAMERA_FRAMES_TO_SKIP,
        jpeg_qual=Constants.JPEG_QUALITY,
        cam_name=Constants.USB_CAMERA_NAME,
        time=now.strftime("%Y.%m.%d %H:%M:%S"),
        dir=Constants.SNAPSHOT_DIR,
        file=image_filename))

    logging.info("Uploading image to Google Drive.")
    currentYear = now.strftime("%Y")
    currentMonth = now.strftime("%b")
    currentDay = now.strftime("%d")
    folder_id = ""
    yearFolderId = ""

    gdrive = gdrive.GDrive()
    gdrive.pushFile(os.path.join(Constants.SNAPSHOT_DIR, image_filename), "Cameras/{year}/{month}/{day}".format(year = currentYear, month = currentMonth, day = currentDay)) 

    time.sleep(Constants.SEC_TO_WAIT_FOR_UPLOAD)

    cleanup()

    logging.info("Total duration: {}".format(datetime.datetime.now() - startTime))
    logging.info("================================= END =================================================")

except Exception, e:
    logging.error("Exception caught: {}".format(str(e)))
    cleanup()
    logging.fatal("Exiting unexpectedly.")
    exit(1)
