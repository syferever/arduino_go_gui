void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);
}

void loop() {
  // put your main code here, to run repeatedly:
  for (int i = 1; i <=12; i++) {
    Serial.print(String(i) + '\n');
    delay(1000);
  }
}
