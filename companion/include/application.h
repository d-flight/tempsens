#ifndef TEMPSENS_APPLICATION_H
#define TEMPSENS_APPLICATION_H

#include "PubSubClient.h"

#include "drivers.h"
#include "config.h"
#include "controller.h"

namespace tempsens {
    class Application {
    private:
        config::Wifi* wifiConfig;
        config::Gpio* gpioConfig;
        config::I2C* i2cConfig;
        PubSubClient* mqttClient;
        drivers::RgbLed* statusLed;
        Controller* controller;

        void setupMqtt(const config::Mqtt& mqttConfig);
        void setupWifi();

    public:
        Application(
            config::Wifi* wifiConfig,
            config::Mqtt* mqttConfig, 
            config::Gpio* gpioConfig,
            config::I2C* i2cConfig
        );

        Controller* boot();

        void connect(bool reconnect);
    };

}

#endif