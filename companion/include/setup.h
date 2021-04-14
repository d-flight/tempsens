
// GPIO
#define TEMPSENS_PIN_RGB_R 12
#define TEMPSENS_PIN_RGB_G 15
#define TEMPSENS_PIN_RGB_B 13

#define TEMPSENS_PIN_RELAY 5

#define TEMPSENS_PIN_I2C_SCL 0
#define TEMPSENS_PIN_I2C_SDA 4

// I2C
#define TEMPSENS_I2C_ADDRESS_BME280 0x76

// WiFi
#define TEMPSENS_WIFI_HOSTNAME "tempsens-companion" // also used for mqtt
#define TEMPSENS_WIFI_SSID "my-wifi-ap"
#define TEMPSENS_WIFI_PWD "password123"

// MQTT
#define TEMPSENS_MQTT_SERVER "mosquitto.local"
#define TEMPSENS_MQTT_PORT 1883

// App
#define TEMPSENS_SLEEP_INTERVAL 10 * 1000 // in ms