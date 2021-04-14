#ifndef TEMPSENS_CONFIG_H
#define TEMPSENS_CONFIG_H

#include <Arduino.h>


namespace tempsens { namespace config {

    struct Wifi {
        std::string hostname;
        std::string ssid;
        std::string password;

        Wifi(std::string hostname, std::string ssid, std::string password);
    };

    struct Mqtt {
        std::string server;
        int port;

        Mqtt(std::string server, int port);
    };

    struct Gpio {
        unsigned char pinLedR,
            pinLedG,
            pinLedB,

            pinRelay,

            pinI2cClock,
            pinI2cData;

        Gpio(
            unsigned char pinLedR, 
            unsigned char pinLedG, 
            unsigned char pinLedB, 
            unsigned char pinRelay, 
            unsigned char pinI2cClock, 
            unsigned char pinI2cData
        );
    };

    struct I2C {
        unsigned char bme280Address;

        I2C(unsigned char bme280Address);
    };

}}
#endif