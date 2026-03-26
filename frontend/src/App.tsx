import "./App.css";
import Footer from "./components/Footer";
import Map from "./components/Map";
import Navbar from "./components/Navbar";
import Room from "./components/Room";

function App({ roomSlug }: { roomSlug: string }) {
  return (
    <>
      <div className="app-shell">
        <Navbar />
        {roomSlug == "/" ? <Map /> : <Room roomSlug={roomSlug} />}
      </div>
      <Footer />
    </>
  );
}

export default App;
