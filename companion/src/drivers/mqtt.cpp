#include <PubSubClient.h>
#include <ArduinoJson.h>

#include "drivers.h"

#define TEMPSENS_MQTT_TOPIC_REPORT "tempsens/report"
#define TEMPSENS_MQTT_TOPIC_CONTROL "tempsens/control"

#define TEMPSENS_CONTROL_TEMPERATURE 1
#define TEMPSENS_CONTROL_TOGGLE_ACTIVE 2

using namespace tempsens::drivers;


void invalidMessage(byte* message) {
    Serial.print("Received invalid message: ");
    Serial.println(*message);
};

void Mqtt::onMqttMessage(char* topic, byte* payload, unsigned int length) {
    if (std::string(TEMPSENS_MQTT_TOPIC_CONTROL) != std::string(topic)) {
        return invalidMessage(payload);
    }

    // decode the json
    StaticJsonDocument<64> message;
    deserializeJson(message, payload);

    // handle command
    auto type = message["Type"].as<char>();

    switch (type) {
        case TEMPSENS_CONTROL_TEMPERATURE: {
            short desiredTemperature = message["Desired"].as<short>();
            if (0 == desiredTemperature || !data::isValidTemperature(desiredTemperature)) {
                return invalidMessage(payload);
            }
            this->changeDesiredTemperatureMessageHandler(desiredTemperature);
        }
        break;

        case TEMPSENS_CONTROL_TOGGLE_ACTIVE: {
            bool newState = message["Active"].as<bool>();
            this->toggleActiveMessageHandler(newState);
        }
        break;

        default:
            return invalidMessage(payload);
    }
};

Mqtt::Mqtt(PubSubClient* mqttClient) : mqttClient(mqttClient) {
    using namespace std::placeholders;
    mqttClient->setCallback([this] (char* topic, byte* payload, unsigned int length) {
        this->onMqttMessage(topic, payload, length);
        Serial.println("Received mqtt message");
    });

    auto subscribed = this->mqttClient->subscribe(TEMPSENS_MQTT_TOPIC_CONTROL, 1); // at least once
    Serial.printf("subscribed: %s\n", (subscribed ? "yes" : "no"));
};

bool Mqtt::report(data::Reading reading, int desiredTemperature, data::HeatingState heatingState) {
    // build json
    StaticJsonDocument<32> readingDocument;
    readingDocument["Temperature"] = reading.temperature;
    readingDocument["Humidity"] = reading.humidity;
    
    StaticJsonDocument<128> rootDocument;
    rootDocument["Desired"] = desiredTemperature;
    rootDocument["Reading"] = readingDocument;
    rootDocument["HeatingState"] = heatingState;

    // serialize
    String message = "";
    serializeJson(rootDocument, message);

    Serial.print("temperature: ");
    Serial.println(reading.temperature);
    Serial.print("humidity: ");
    Serial.println(reading.humidity);

    // ship
    return this->mqttClient->publish(TEMPSENS_MQTT_TOPIC_REPORT, message.c_str());
};

void Mqtt::onChangeDesiredTemperatureMessage(std::function<void (signed short)> callback) {
    this->changeDesiredTemperatureMessageHandler = callback;
};

void Mqtt::onToggleActiveMessage(std::function<void (bool toggle)> callback) {
    this->toggleActiveMessageHandler = callback;
};