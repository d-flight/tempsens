#include <Arduino.h>
#include <ESP8266WiFi.h>
#include <PubSubClient.h>

#include "config.h"
#include "drivers.h"
#include "application.h"
#include "controller.h"

#include "setup.h"

using namespace tempsens;

Controller* controller;
Application* app;

void setup() {
    // setup pins
    auto gpioConfig = new config::Gpio(
        TEMPSENS_PIN_RGB_R,
        TEMPSENS_PIN_RGB_G,
        TEMPSENS_PIN_RGB_B,
        TEMPSENS_PIN_RELAY,
        TEMPSENS_PIN_I2C_SCL,
        TEMPSENS_PIN_I2C_SDA
    );

    // setup i2c
    auto i2cConfig = new config::I2C(TEMPSENS_I2C_ADDRESS_BME280);

    // setup wifi
    auto wifiConfig = new config::Wifi(TEMPSENS_WIFI_HOSTNAME, TEMPSENS_WIFI_SSID, TEMPSENS_WIFI_PWD);

    // setup mqtt
    auto mqttConfig = new config::Mqtt(TEMPSENS_MQTT_SERVER, TEMPSENS_MQTT_PORT);
    
    // setup application
    app = new Application(wifiConfig, mqttConfig, gpioConfig, i2cConfig);

    // boot
    controller = app->boot();
}

void loop() {
    app->reconnect();

    controller->tick();

    delay(TEMPSENS_SLEEP_INTERVAL);
}