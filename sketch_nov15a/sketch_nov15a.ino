bool stringComplete = false;
String inputString = "";
int LED = LED_BUILTIN;

void SerialPrint(String s) {
  Serial.print(s + "\n");
}

void setup() {
  // put your setup code here, to run once:
  pinMode(LED_BUILTIN, OUTPUT);
  Serial.begin(115200);
}

void loop() {
  // put your main code here, to run repeatedly:
  if (stringComplete) {
    char command = inputString.charAt(0);
    int argument = inputString.substring(1).toInt();
    // Serial.println(argument);
    switch (command) {
      case 'l':
        digitalWrite(LED, argument);
        SerialPrint(argument ? "Light on" : "Light off");
        break;
      case 'm':
        digitalWrite(LED, 1);
        delay(argument*1000);
        digitalWrite(LED, 0);
        break;
      case 'd':
        for (int i = 1; i <= 10; i++) {
          SerialPrint(String(i));
        }
        break;
      default:
        SerialPrint("unknown command: " + String(command));
    }
    inputString = "";
    stringComplete = false;
  }
}

void serialEvent() {
  while (Serial.available()) {
    char inChar = (char)Serial.read();
    if (inChar == '\n') {
      stringComplete = true;
    } else {
      inputString += inChar;
    }
  }
}
