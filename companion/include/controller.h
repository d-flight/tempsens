#ifndef TEMPSENS_CONTROLLER_H
#define TEMPSENS_CONTROLLER_H

#include <PubSubClient.h>

#include "drivers.h"
#include "data.h"

namespace tempsens {

    class Controller {
        private:
            drivers::Mqtt* mqtt;
            drivers::Bme280* bme;
            drivers::Relay* relay;
            drivers::RgbLed* led;

            int desiredTemperature = 2350;
            data::HeatingState heatingState = data::HeatingState::off;

            void updateStatusLed();
            void updateHeatingState(data::Reading reading);
            void setDesiredTemperature(signed short temperature);

        public:
            Controller(drivers::Mqtt* mqtt, drivers::Bme280* bme, drivers::Relay* relay, drivers::RgbLed* led);

            void tick();
            void onReconnect();
    };


};


#endif