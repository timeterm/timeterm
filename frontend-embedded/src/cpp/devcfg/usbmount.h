#pragma once

#include <QString>
#include <optional>

std::optional<QString> tryMountConfig();
std::optional<QString> tryUnmountConfig();
