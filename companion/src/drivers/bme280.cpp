#include <Wire.h>

#include "drivers.h"
#include "data.h"

using namespace tempsens::drivers;

Bme280::Bme280(unsigned char address, unsigned char sdaPin, unsigned char sclPin) {
    // setup i2c
    Wire.pins(sdaPin, sclPin);

    Serial.println("Wire setup");

    // setup bme
    this->bme = Adafruit_BME280();
    Serial.println("bme setup");

    bool success = false;
    
    while (!this->bme.begin(address, &Wire)) {
        Serial.println("connecting..");
        delay(1000);
    }

    if (success) {
        Serial.println("bme begin successful");
    } else {
        Serial.println("failed to begin reading bme");
    }
};

tempsens::data::Reading Bme280::getLatestReading() {
    return tempsens::data::Reading(
        this->bme.readTemperature() * 100,
        this->bme.readHumidity() * 100
    );
};

