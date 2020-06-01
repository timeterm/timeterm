#include "cardreadercontroller.h"

CardReaderController::CardReaderController(CardReader *cardReader,
                                           QObject *parent)
    : QObject(parent), m_cardReader(cardReader)
{
    m_cardReader->moveToThread(&cardReaderThread);
    connect(&cardReaderThread, &QThread::finished, m_cardReader, &QObject::deleteLater);
    connect(m_cardReader, &CardReader::cardRead, this, &CardReaderController::cardRead);
    cardReaderThread.start();
}

CardReaderController::~CardReaderController() {
    m_cardReader->shutDown();
    cardReaderThread.quit();
    cardReaderThread.wait();
}
