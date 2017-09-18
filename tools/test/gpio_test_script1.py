#!/usr/bin/python
import RPi.GPIO as GPIO
import time

GPIO.setmode(GPIO.BCM)

# init list with pin numbers

pinList = [26, 16, 20, 21, 5, 6, 13, 19]


# loop through pins and set mode and state to 'low'

for i in pinList:
    GPIO.setup(i, GPIO.OUT)
    GPIO.output(i, GPIO.LOW)

# time to sleep between operations in the main loop

SleepTimeL = 2

# main loop

try:
  for i in pinList:
      print "----> {}".format(i)
      GPIO.output(i, GPIO.HIGH)
      time.sleep(SleepTimeL);
      
  GPIO.cleanup()
  print "Good bye!"

# End program cleanly with keyboard
except KeyboardInterrupt:
  print "  Quit"

  # Reset GPIO settings
  GPIO.cleanup()


# find more information on this script at
# http://youtu.be/oaf_zQcrg7g
