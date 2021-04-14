#ifndef TEMPSENS_DATA_H
#define TEMPSENS_DATA_H

namespace tempsens { namespace data {

    struct Reading {
        short temperature;
        short humidity;

        Reading(short temperature, short humidity);
    };

    enum HeatingState { 
        off = 0,
        idle = 1,
        heating =2 
    };


    bool isValidTemperature(short temp);
}}

#endif