export type Machine = {
  id: number; // In room ID, unique within the room
  is_operational: boolean;
  is_washer: boolean;
  end_at: number | null; // Unix timestamp in milliseconds
};

export type RoomData = {
  washers: Machine[];
  dryers: Machine[];
};

export type MachineStatus = {
  isAvailable: boolean;
  label: string;
  timer: string | null;
};

function getRemainingSeconds(machine: Machine, nowMs = Date.now()): number {
  if (!machine.end_at) {
    return 0;
  }

  const diffSeconds = (machine.end_at - nowMs) / 1000;

  return diffSeconds > 0 ? Math.floor(diffSeconds) : 0;
}

export function getMachineStatus(
  machine: Machine,
  nowMs = Date.now(),
): MachineStatus {
  if (!machine.is_operational) {
    return {
      isAvailable: false,
      label: "Out of Order",
      timer: null,
    };
  }

  const totalSeconds = getRemainingSeconds(machine, nowMs);
  if (totalSeconds === 0) {
    return {
      isAvailable: true,
      label: "Available",
      timer: null,
    };
  }

  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;

  return {
    isAvailable: false,
    label: "In Use",
    timer: `${minutes}:${String(seconds).padStart(2, "0")}`,
  };
}
