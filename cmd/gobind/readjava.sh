#!/bin/sh
unzip iot.aar
jar xf classes.jar
javap -c iot/Iot.class
