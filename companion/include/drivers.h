#ifndef TEMPSENS_DRIVERS_H
#define TEMPSENS_DRIVERS_H

#include <PubSubClient.h>
#include <Adafruit_BME280.h>

#include "data.h"


namespace tempsens { namespace drivers {

    class RgbLed {
        private:
            unsigned char pinR, powerR,
                pinG, powerG,
                pinB, powerB;

            void propagate();

        public:
            RgbLed(unsigned char rPin, unsigned char gPin, unsigned char bPin);

            void display(unsigned char r, unsigned char g, unsigned char b);
            void blink(unsigned char r, unsigned char g, unsigned char b, unsigned long timeOn = 500, unsigned long timeOff = 200);
            void off();
    };

    class Relay {
        private:
            unsigned char pin;
            bool power = false;

            void propagate();

        public:
            Relay(unsigned char pin);

            void flip();
            void toggle(bool toggle);
    };

    class Bme280 {
        private:
            Adafruit_BME280 bme;

        public:
            Bme280(unsigned char address, unsigned char sdaPin, unsigned char sclPin);

            data::Reading getLatestReading();
    };

    class Mqtt {
        private:
            PubSubClient* mqttClient;

            std::function<void (short)> changeDesiredTemperatureMessageHandler;
            std::function<void (bool)> toggleActiveMessageHandler;

            void onMqttMessage(char* topic, byte* payload, unsigned int length);
        public:
            Mqtt(PubSubClient* mqttClient);

            bool report(data::Reading reading, int desiredTemperature, data::HeatingState heatingState);
            void onChangeDesiredTemperatureMessage(std::function<void (short)> callback);
            void onToggleActiveMessage(std::function<void (bool toggle)> callback);
    };

}}


#endif