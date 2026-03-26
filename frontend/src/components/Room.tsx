import useSWR from "swr";
import MachineButton from "./MachineButton";
import type { Machine, RoomData } from "../utils";
const apiUrl = import.meta.env.VITE_API_URL ?? "";

function slugToName(slug: string) {
  return slug
    .split("-")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}

// GET /api/rooms/{room_slug}/machines
// {
//     "washers": [
//         {
//             "id": 1,
//             "is_operational": true,
//             "is_washer": true,
//             "end_at": null
//         }
//     ],
//     "dryers": [
//         {
//             "id": 7,
//             "is_operational": true,
//             "is_washer": false,
//             "end_at": null
//         }
//     ]
// }

const fetcher = (arg: string) => {
  console.log("Fetching", arg);
  return fetch(arg).then((res) => res.json());
};
function Room({ roomSlug }: { roomSlug: string }) {
  const { data, error, isLoading, mutate } = useSWR<RoomData>(
    `${apiUrl}/api/rooms/${roomSlug}/machines`,
    fetcher,
    {
      refreshInterval: 15000, // Refresh every 15 seconds
    },
  );

  if (error) {
    console.log(error);
    return (
      <div className="status-container-centered">
        <div className="error-inline">
          <p className="error-title">Error loading room data</p>
          <p className="error-subtext">Please try again later.</p>
        </div>
      </div>
    );
  }

  if (isLoading || !data) {
    console.log("Loading...");
    return (
      <div className="status-container-centered">
        <div className="spinner" />
      </div>
    );
  }

  return (
    <>
      <main className="room-page">
        <h1>{slugToName(roomSlug)} Laundry Room</h1>
        <section className="machine-section" aria-label="Washing Machines">
          <h2 className="machine-section-title">Washing Machines</h2>
          <div className="machine-grid">
            {data.washers.map((machine: Machine) => (
              <MachineButton
                key={machine.id}
                machine={machine}
                roomSlug={roomSlug}
                onMachineUpdated={() => {
                  void mutate();
                }}
              />
            ))}
          </div>
        </section>
        <section className="machine-section" aria-label="Drying Machines">
          <h2 className="machine-section-title">Drying Machines</h2>
          <div className="machine-grid">
            {data.dryers.map((machine: Machine) => (
              <MachineButton
                key={machine.id}
                machine={machine}
                roomSlug={roomSlug}
                onMachineUpdated={() => {
                  void mutate();
                }}
              />
            ))}
          </div>
        </section>
      </main>
    </>
  );
}

export default Room;
