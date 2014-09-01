package ioprio

import "syscall"

// taken from <linux/ioprio.h>
// scheduling classes
type Class uint
const (
    None       = 0
    RealTime   = 1
    BestEffort = 2
    Idle       = 3
)

func (c Class) String() string {
    switch c {
    case None:
        return "none"
    case RealTime:
        return "realtime"
    case BestEffort:
        return "best-effort"
    case Idle:
        return "idle"
    default:
        return "unknown"
    }
}

// which
type Which uint
const (
    Process      = 1
    ProcessGroup = 2
    User         = 3
)

func (w Which) String() string {
    switch w {
    case Process:
        return "process"
    case ProcessGroup:
        return "process group"
    case User:
        return "user"
    default:
        panic("invalid")
    }
}

type Prio uint

const Normal = 4
const BestEffortNr = 8

const ioPrioBits = 16
const ioPrioClassShift = 13
const ioPrioPrioMask = ((1 << ioPrioClassShift) - 1)

func ioPrioClass(mask uint) Class {
    return Class(mask >> ioPrioClassShift)
}

func ioPrioData(mask uint) Prio {
    return Prio(mask & ioPrioPrioMask)
}

func ioPrioValue(class Class, data Prio) uint {
    return (uint(class) << ioPrioClassShift) | uint(data)
}

func SetIoPrio(which Which, who uint, class Class, prio Prio) error {
    _, _, err := syscall.Syscall(syscall.SYS_IOPRIO_SET, uintptr(which), uintptr(who), uintptr(ioPrioValue(class, prio)))
    return err
}

// returns class & prio
func GetIoPrio(which Which, who uint) (Class, Prio, error) {
    prio, _, err := syscall.Syscall(syscall.SYS_IOPRIO_GET, uintptr(which), uintptr(who), uintptr(0))
    if err != 0 {
        return Class(0), Prio(0), err
    } else {
        class := ioPrioClass(uint(prio))
        value := ioPrioData(uint(prio))
        return Class(class), Prio(value), nil
    }
}
