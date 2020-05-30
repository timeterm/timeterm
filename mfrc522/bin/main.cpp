#include <iostream>

#include <csignal>
#include <mfrc522/mfrc522.h>

volatile bool quit = false;

void handleInterrupt(int)
{
    quit = true;
}

int main(int argc, char *argv[])
{
    auto mfrc522 = Mfrc522::Device{};
    mfrc522.init();

    std::signal(SIGINT, handleInterrupt);

    while (!quit) {
        auto [status, tagType] = mfrc522.request(Mfrc522::PiccCommand::ReqIdl);

        if (status != Mfrc522::Status::Ok) {
            continue;
        }
        std::cout << "Card detected" << std::endl;

        auto [newStatus, uid] = mfrc522.antiColl();
        if (newStatus != Mfrc522::Status::Ok) {
            continue;
        }

        std::cout << "Card UID: " << std::to_string(uid[0]) << std::to_string(uid[1])
                  << std::to_string(uid[2]) << std::to_string(uid[3]) << std::endl;
    }
}