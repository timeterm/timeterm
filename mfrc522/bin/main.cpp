
#ifdef WIN32
#include <windows.h>
#else
#include <mfrc522/mfrc522.h>
#include <unistd.h>
#endif

void delay(int ms){
#ifdef WIN32
    Sleep(ms);
#else
    usleep(ms*1000);
#endif
}

int main(){
    MFRC522 mfrc;

    mfrc.PCD_Init();

    while(1){
        // Look for a card
        if(!mfrc.PICC_IsNewCardPresent())
            continue;

        if( !mfrc.PICC_ReadCardSerial())
            continue;

        // Print UID
        for(uint8_t i = 0; i < mfrc.uid.size; ++i){
            if(mfrc.uid.uidByte[i] < 0x10){
                printf(" 0");
                printf("%X",mfrc.uid.uidByte[i]);
            }
            else{
                printf(" ");
                printf("%X", mfrc.uid.uidByte[i]);
            }

        }
        printf("\n");
        delay(1000);
    }
    return 0;
}
