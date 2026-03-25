import type { Machine } from "../utils";
import washerImg from "../assets/washer.svg";
import drierImg from "../assets/drier.svg";
import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import MachineModal from "./MachineModal";
import {
  getMachineStatusLabel,
  getMachineTimer,
  isMachineAvailable,
} from "../utils";
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
  const [nowSec, setNowSec] = useState(Math.floor(Date.now() / 1000));

  useEffect(() => {
    const intervalId = window.setInterval(() => {
      setNowSec(Math.floor(Date.now() / 1000));
    }, 1000);

    return () => {
      window.clearInterval(intervalId);
    };
  }, []);

  const isAvailable = isMachineAvailable(machine, nowSec);
  const status = getMachineStatusLabel(machine, nowSec);
  const timer = getMachineTimer(machine, nowSec);
  return (
    <>
      <button
        key={machine.id}
        className={clsx("machine-button", { "not-available": !isAvailable })}
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
        <span className="machine-status">{timer ? timer : status}</span>
      </button>
      {showModal &&
        createPortal(
          <MachineModal
            machine={machine}
            roomSlug={roomSlug}
            onMachineUpdated={onMachineUpdated}
            onClose={() => setShowModal(false)}
          />,
          document.body,
        )}
    </>
  );
}

export default MachineButton;
