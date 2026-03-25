import { getMachineStatus, type Machine } from "../utils";
import washerImg from "../assets/washer.svg";
import drierImg from "../assets/drier.svg";
import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import MachineModal from "./MachineModal";
import { clsx } from "clsx";

type MachineButtonProps = {
  machine: Machine;
  roomSlug: string;
  onMachineUpdated: () => void;
};

function MachineButton({
  machine,
  roomSlug,
  onMachineUpdated,
}: MachineButtonProps) {
  const [showModal, setShowModal] = useState(false);
  const [nowMs, setNowMs] = useState(Date.now());

  useEffect(() => {
    const intervalId = window.setInterval(() => {
      setNowMs(Date.now());
    }, 1000);

    return () => {
      window.clearInterval(intervalId);
    };
  }, []);

  const machineStatus = getMachineStatus(machine, nowMs);

  return (
    <>
      <button
        key={machine.id}
        className={clsx("machine-button", {
          "not-available": !machineStatus.isAvailable,
        })}
        onClick={() => setShowModal(true)}
      >
        <img
          src={machine.is_washer ? washerImg : drierImg}
          alt=""
          aria-hidden="true"
          className="machine-icon"
          width={400}
          height={400}
        />
        <span className="machine-number">#{machine.id}</span>
        <span className="machine-status">
          {machineStatus.timer ? machineStatus.timer : machineStatus.label}
        </span>
      </button>
      {showModal &&
        createPortal(
          <MachineModal
            machine={machine}
            roomSlug={roomSlug}
            onMachineUpdated={onMachineUpdated}
            onClose={() => setShowModal(false)}
            machineStatus={machineStatus}
          />,
          document.body,
        )}
    </>
  );
}

export default MachineButton;
