
export type Machine = {
  id: number;
  is_operational: boolean;
  is_washer: boolean;
  end_at: number | null;
};

export type RoomData = {
  washers: Machine[];
  dryers: Machine[];
};

type MachineStatus = {
  isAvailable: boolean;
  status: string;
  timer: string | null;
};

function getNowSec(): number {
  return Math.floor(Date.now() / 1000);
}

function getRemainingSeconds(machine: Machine, nowSec: number): number {
  if (!machine.end_at) {
    return 0;
  }

  const diffSeconds = machine.end_at - nowSec;

  return diffSeconds > 0 ? diffSeconds : 0;
}

export function isMachineAvailable(machine: Machine, nowSec = getNowSec()): boolean {
  if (!machine.is_operational) {
    return false;
  }

  return getRemainingSeconds(machine, nowSec) === 0;
}

export function getMachineStatusLabel(
  machine: Machine,
  nowSec = getNowSec(),
): "Available" | "In Use" | "Out of Order" {
  if (!machine.is_operational) {
    return "Out of Order";
  }

  return getRemainingSeconds(machine, nowSec) > 0 ? "In Use" : "Available";
}

export function getMachineTimer(machine: Machine, nowSec = getNowSec()): string | null {
  const totalSeconds = getRemainingSeconds(machine, nowSec);

  if (totalSeconds === 0) {
    return null;
  }

  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;

  return `${minutes}:${String(seconds).padStart(2, "0")}`;
}

export function getMachineStatus(machine: Machine, nowSec = getNowSec()): MachineStatus {
  const status = getMachineStatusLabel(machine, nowSec);
  const timer = getMachineTimer(machine, nowSec);

  return {
    isAvailable: status === "Available",
    status,
    timer,
  };
}
