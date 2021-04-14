#include "config.h"

using namespace tempsens::config;

Wifi::Wifi(std::string hostname, std::string ssid, std::string password) {
    this->hostname = hostname;
    this->ssid = ssid;
    this->password = password;
};

Mqtt::Mqtt(std::string server, int port) {
    this->server = server;
    this->port = port;
};

Gpio::Gpio(unsigned char pinLedR, unsigned char pinLedG, unsigned char pinLedB, unsigned char pinRelay, unsigned char pinI2cClock, unsigned char pinI2cData):
    pinLedR(pinLedR),
    pinLedG(pinLedG),
    pinLedB(pinLedB),
    pinRelay(pinRelay),
    pinI2cClock(pinI2cClock),
    pinI2cData(pinI2cData) {};

I2C::I2C(unsigned char bme280Address): bme280Address(bme280Address) {};