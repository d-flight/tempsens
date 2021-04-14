#include "data.h"

using namespace tempsens::data;

Reading::Reading(short temperature, short humidity) :
        temperature(temperature), humidity(humidity) {};

bool tempsens::data::isValidTemperature(short temp) {
    return -4000 <= temp && temp <= 8500;
};
