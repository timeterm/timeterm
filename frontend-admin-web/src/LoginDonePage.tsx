import React from "react";
import { Typography } from "@rmwc/typography";
import "@rmwc/typography/styles";
import { Theme } from "@rmwc/theme";
import { Elevation } from "@rmwc/elevation";
import Logo from "./logo-black.svg";
import { useLocation, useHistory } from "react-router-dom";
import Cookies from "universal-cookie";

const useQuery = () => {
  return new URLSearchParams(useLocation().search);
};

const LoginDonePage: React.FC = (props) => {
  const query = useQuery();
  const history = useHistory();

  const status = query.get("status");
  const token = query.get("token");
  if (status === "ok" && token) {
    const tokenDec = atob(token);
    new Cookies().set("ttsess", tokenDec, {
      secure: true,
      path: "/",
    });
    history.push("/devices");
  }
  const error = query.get("error");

  return (
    <div
      className={"LoginPage"}
      style={{
        width: "100%",
        height: "100%",
        backgroundSize: "400% 400%",
        background:
          "#efefef linear-gradient(120deg, #7fc1ff 10%, #2a73b4 100%)",
        transition: "0.3s",
      }}
    >
      <Theme use={["surface"]} wrap>
        <Elevation
          z={16}
          style={{
            padding: 32,
            height: "100%",
            boxSizing: "border-box",
            width: "26em",
            justifyContent: "space-between",
            display: "flex",
            flexDirection: "column",
          }}
        >
          <div style={{ display: "flex", flexDirection: "column" }}>
            <img
              src={Logo}
              alt={"Timeterm Logo"}
              style={{ width: 96, marginBottom: 16 }}
            />
            <Typography use="headline4" style={{ marginBottom: 8 }}>
              {status === "ok" ? "Ingelogd" : "Er is een fout opgetreden"}
            </Typography>

            <span>{status !== "ok" && error}</span>
          </div>

          <span>&copy; {new Date().getFullYear()} de auteurs van Timeterm</span>
        </Elevation>
      </Theme>
    </div>
  );
};

export default LoginDonePage;
