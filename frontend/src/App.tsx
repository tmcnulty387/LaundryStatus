import "./App.css";
import Map from "./components/Map";
import Navbar from "./components/Navbar";
import Room from "./components/Room";

function App({ roomSlug }: { roomSlug: string }) {
  return (
    <>
      <Navbar />
      {roomSlug == "/" ? <Map /> : <Room roomSlug={roomSlug} />}
    </>
  );
}

export default App;
