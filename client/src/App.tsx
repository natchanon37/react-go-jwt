import React, { useEffect, useState } from "react";
import { Navbar } from "./components/nav";
import { Login } from "./pages/Login";
import { Register } from "./pages/Register";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Home } from "./pages/Home";
import "./App.css";

function App() {
  const [name, setName] = useState("");
  useEffect(() => {
    const fetchUser = async () => {
      const data = await fetch("http://127.0.0.1:8000/api/user", {
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      });
      const content = await data.json();
      console.log(content.name);
      setName(content.name);
    };

    fetchUser();
  }, [name]);
  return (
    <>
      <BrowserRouter>
        <Navbar name={name} setName={setName} />
        <Routes>
          <Route path="/" element={<Home name={name} />} />
          <Route path="/register" element={<Register />} />
          <Route path="/login" element={<Login setName={setName} />} />
        </Routes>
      </BrowserRouter>
    </>
  );
}

export default App;
