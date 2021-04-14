#include <Arduino.h>

#include "drivers.h"

using namespace tempsens::drivers;

void RgbLed::propagate() {
    analogWrite(pinR, powerR);
    analogWrite(pinG, powerG);
    analogWrite(pinB, powerB);
}

RgbLed::RgbLed(unsigned char rPin, unsigned char gPin, unsigned char bPin) {
    pinR = rPin;
    pinG = gPin;
    pinB = bPin;

    pinMode(pinR, OUTPUT);
    pinMode(pinG, OUTPUT);
    pinMode(pinB, OUTPUT);
}

void RgbLed::display(unsigned char r, unsigned char g, unsigned char b) {
    powerR = r;
    powerG = g;
    powerB = b;

    propagate();
}

void RgbLed::blink(unsigned char r, unsigned char g, unsigned char b, unsigned long timeOn, unsigned long timeOff) {
    this->display(r, g, b);
    delay(timeOn);

    this->off();
    delay(timeOff);
}

void RgbLed::off() {
    this->display(0, 0, 0);
};