#include "controller.h"

using namespace tempsens;

#define TEMPSENS_CONTROLLER_UPPER_TEMP_BUFFER 25
#define TEMPSENS_CONTROLLER_LOWER_TEMP_BUFFER 15

Controller::Controller(drivers::Mqtt* mqtt, drivers::Bme280* bme, drivers::Relay* relay, drivers::RgbLed* led):
    mqtt(mqtt),
    bme(bme),
    relay(relay),
    led(led),
    heatingState(data::HeatingState::idle) {

        // listen for commands
        this->mqtt->onChangeDesiredTemperatureMessage([this](short temperature) {
            this->desiredTemperature = temperature;
            Serial.printf("Changed desired temperature to: %d\n", temperature);
        });

        this->mqtt->onToggleActiveMessage([this](bool toggle) {
            if (false == toggle) { // turned off
                this->heatingState = data::HeatingState::off;
                Serial.println("Switched heating state to off");
            
            // was turned off before, now switch to idle. next tick() will pick it up
            } else if (data::HeatingState::off == this->heatingState) {
                this->heatingState = data::HeatingState::idle;
                Serial.println("Switched heating state to idle");
            }
        });
};

void Controller::updateHeatingState(data::Reading reading) {
    // do nothing when turned off
    if (data::HeatingState::off == this->heatingState) {
        return;
    }

    // temperature higher than buffer, go to idle
    if (TEMPSENS_CONTROLLER_UPPER_TEMP_BUFFER < reading.temperature - this->desiredTemperature) {
        this->heatingState = data::HeatingState::idle;
    
    // temperature lower than buffer, go to heating
    } else if (TEMPSENS_CONTROLLER_LOWER_TEMP_BUFFER < this->desiredTemperature - reading.temperature) {
        this->heatingState = data::HeatingState::heating;
    }
};

void Controller::tick() {
    // fetch the latest reading
    data::Reading reading = this->bme->getLatestReading();

    // calculate heating state
    this->updateHeatingState(reading);

    // update relay
    this->relay->toggle(data::HeatingState::heating == this->heatingState);

    // update status led
    this->updateStatusLed();

    // publish report
    auto success = this->mqtt->report(reading, this->desiredTemperature, this->heatingState);

    if (success) {
        Serial.println("successfully published mqtt message");
    } else {
        Serial.println("error when trying to publish mqtt message");
    }
};

void Controller::updateStatusLed() {
    switch (this->heatingState) {
        case data::HeatingState::off: 
            this->led->off();
            break;
        case data::HeatingState::idle:
            this->led->display(0, 0, 0xff);
            break;
        case data::HeatingState::heating:
            this->led->display(0xff, 0xa5, 0);
            break;
    }
};