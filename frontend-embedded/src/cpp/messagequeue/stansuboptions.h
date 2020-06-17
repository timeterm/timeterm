#ifndef STANSUBOPTIONS_H
#define STANSUBOPTIONS_H

#include "scopedpointer.h"

#include <nats.h>

#include <QObject>

namespace MessageQueue
{

using StanSubOptionsDeleter = ScopedPointerDestroyerDeleter<stanSubOptions, void, &stanSubOptions_Destroy>;
using StanSubOptionsScopedPointer = QScopedPointer<stanSubOptions, StanSubOptionsDeleter>;

class StanSubOptions: public QObject
{
    Q_OBJECT

public:
    explicit StanSubOptions(QObject *parent = nullptr);

private:
    StanSubOptionsScopedPointer m_subOptions;
};

}

#endif // STANSUBOPTIONS_H
