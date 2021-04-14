#include <ESP8266WiFi.h>
#include <PubSubClient.h>
#include <Wire.h>

#include "setup.h"
#include "application.h"

using namespace tempsens;

void Application::reconnect() {
    // fixing lost wifi connection
    while (WiFi.status() != WL_CONNECTED) {
        this->statusLed->blink(0x0, 0x0, 0xff);
        Serial.println("connecting to wifi..");
    }

    Serial.println("wifi connected.");

    // fixing lost mqtt connection
    while (!this->mqttClient->connected()) {
        Serial.println("Connecting to MQTT...");
        this->statusLed->blink(0xff, 0x0, 0x0);

    
        if (!this->mqttClient->connect(this->wifiConfig->hostname.c_str())) {
            Serial.print("failed with state ");
            Serial.print(this->mqttClient->state());
        }
    }

    this->mqttClient->loop();

    Serial.println("mqtt connected");
};

void Application::setupWifi() {
    WiFi.hostname(this->wifiConfig->hostname.c_str());
    WiFi.setAutoReconnect(true);
    WiFi.begin(this->wifiConfig->ssid.c_str(), this->wifiConfig->password.c_str());
};

Controller* Application::boot() {
    // prepare wifi
    this->setupWifi();

    // and try to connect
    this->reconnect();

    // setup the controller with the rest of the drivers
    return new Controller(
        new drivers::Mqtt(this->mqttClient),
        new drivers::Bme280(this->i2cConfig->bme280Address, this->gpioConfig->pinI2cData, this->gpioConfig->pinI2cClock),
        new drivers::Relay(this->gpioConfig->pinRelay),
        this->statusLed
    );
};

void Application::setupMqtt(const config::Mqtt& mqttConfig) {
    WiFiClient* wifi = new WiFiClient();

    this->mqttClient = new PubSubClient(*wifi);
    this->mqttClient->setServer(mqttConfig.server.c_str(), mqttConfig.port);
    this->mqttClient->setKeepAlive((TEMPSENS_SLEEP_INTERVAL / 1000) * 3); // 3x the sleep interval
};

Application::Application(
    config::Wifi* wifiConfig, 
    config::Mqtt* mqttConfig, 
    config::Gpio* gpioConfig,
    config::I2C* i2cConfig
) {
    this->wifiConfig = wifiConfig;
    this->gpioConfig = gpioConfig;
    this->i2cConfig = i2cConfig;

    this->setupMqtt(*mqttConfig);

    this->statusLed = new drivers::RgbLed(
        this->gpioConfig->pinLedR,
        this->gpioConfig->pinLedG,
        this->gpioConfig->pinLedB
    );
};