bool stringComplete = false;
String inputString = "";

void setup() {
  // put your setup code here, to run once:
  pinMode(LED_BUILTIN, OUTPUT);
  Serial.begin(9600);
}

void loop() {
  // put your main code here, to run repeatedly:
  if (stringComplete) {
    char command = inputString.charAt(0);
    int argument = inputString.substring(1).toInt();
    // Serial.println(argument);
    switch (command) {
      case 'l':
        digitalWrite(LED_BUILTIN, argument);
        Serial.println(argument ? "Light on" : "Light off");
        break;
      case 'm':
        digitalWrite(LED_BUILTIN, 1);
        delay(argument);
        digitalWrite(LED_BUILTIN, 0);
      case 'd':
        for (int i = 1; i < 100; i++) {
          Serial.print(String(i) + '\n');
        }
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
