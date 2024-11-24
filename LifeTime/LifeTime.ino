String inputString = "";
bool stringComplete = false;
int LED = 2;
int d = 1000;
long int t = 0;
const int points = 100;
int V[3][points];

void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);
  pinMode(LED, OUTPUT);
  inputString.reserve(200);
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
        t = micros();
        V[3][points] = memset(V, 0, sizeof(V));
        for (int j = 0; j < argument; j++) {
          digitalWrite(LED, LOW);
          for (int i = 0; i < points; i++) {
            V[0][i] += analogRead(A7);
            delayMicroseconds(d);
          }
          digitalWrite(LED, HIGH);
          for (int i = 0; i < points; i++) {
            V[1][i] += analogRead(A7);
            delayMicroseconds(d);
          }
          digitalWrite(LED, LOW);
          for (int i = 0; i < points; i++) {
            V[2][i] += analogRead(A7);
            delayMicroseconds(d);
          }
        }
        t = (micros() - t)/1000;
        break;
      case 'd':
        for (int i = 0; i < points; i++) {
          Serial.println(V[argument][i]);
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
