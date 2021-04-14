#include <Arduino.h>

#include "drivers.h"

using namespace tempsens::drivers;


Relay::Relay(unsigned char pin): pin(pin) {
    pinMode(pin, OUTPUT);
};

void Relay::propagate() {
    digitalWrite(this->pin, this->power ? 1 : 0);
};

void Relay::flip() {
    this->toggle(!this->power);
};

void Relay::toggle(bool toggle) {
    this->power = toggle;
    this->propagate();
};