String inputString = "";
bool stringComplete = false;
int LED = LED_BUILTIN;
int d = 10;
long int t = 0;
const int points = 100;
const char output = A5;
int V[points];
float sigma0;

void setup() {
  // put your setup code here, to run once:
  Serial.begin(115200);
  pinMode(LED, OUTPUT);
  inputString.reserve(10);
}

void loop() {
  // put your main code here, to run repeatedly:
  if (stringComplete) {
    char command = inputString.charAt(0);
    int argument = inputString.substring(1).toInt();
    switch (command) {
      case 'l':
        digitalWrite(LED, argument);
        Serial.println(argument ? "Light on" : "Light off");
        break;
      case 'p':
        d = argument;
        break;
      case 'm':
        memset(V, 0, sizeof(V));
        sigma0 = 1.0 / (1023.0 / analogRead(A5) - 1);
        for (int j = 0; j < argument; j++) {
          digitalWrite(LED, HIGH);
          delay(1000);
          digitalWrite(LED, LOW);
          t = micros();
          for (int i = 0; i < points; i++) {
            V[i] += analogRead(A5)/float(argument);
            delayMicroseconds(d);
          }
          t = micros() - t;
        }
        break;
      case 'd':
        for (int i = 0; i < points; i++) {
          float sigma = 1.0 / (1023.0 / V[i] - 1.0) - sigma0;
          if (i > 0) {
            float tau = sigma * (t / points) / ((1.0 / (1023.0 / V[i - 1] - 1.0) - sigma0) - sigma);
            // Serial.print(tau/1000);
            // Serial.print(",");
            Serial.print(String(tau/1000) + '\n');
            // Serial.println(sigma);
          }
        }
        break;
      case 't':
        Serial.println(t);
        break;
    }
    inputString = "";
    stringComplete = false;
  }
}

void serialEvent() {
  while (Serial.available()) {
    char inChar = (char)Serial.read();
    inputString += inChar;
    if (inChar == '\n') {
      stringComplete = true;
    }
  }
}
