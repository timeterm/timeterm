import React, { useState } from "react";
import { Typography } from "@rmwc/typography";
import "@rmwc/typography/styles";
import { Button } from "@rmwc/button";
import { Theme } from "@rmwc/theme";
import { Elevation } from "@rmwc/elevation";
import Logo from "./logo-black.svg";
import Cookies from "universal-cookie";
import { useHistory } from "react-router-dom";

const LoginPage: React.FC = (props) => {
  const [isHovering, setIsHovering] = useState(false);
  const history = useHistory();
  if (new Cookies().get("ttsess")) {
    history.push("/devices");
  }

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
        backgroundPosition: isHovering ? 300 : undefined,
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
              Inloggen
            </Typography>
            Welkom terug!
            <Button
              raised
              style={{
                backgroundColor: "black",
                color: "white",
                marginTop: 64,
              }}
              className={"LoginBtn"}
              icon={
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 23 23">
                  <path fill="#f35325" d="M1 1h10v10H1z" />
                  <path fill="#81bc06" d="M12 1h10v10H12z" />
                  <path fill="#05a6f0" d="M1 12h10v10H1z" />
                  <path fill="#ffba08" d="M12 12h10v10H12z" />
                </svg>
              }
              onMouseOver={() => setIsHovering(true)}
              onMouseOut={() => setIsHovering(false)}
              onClick={() => {
                window.location.href = `https://api.timeterm.nl/oidc/login/microsoft?redirectTo=${encodeURIComponent(
                  window.location.href + "login/done"
                )}`;
              }}
            >
              Inloggen met Microsoft
            </Button>
          </div>

          <span>&copy; {new Date().getFullYear()} de auteurs van Timeterm</span>
        </Elevation>
      </Theme>
    </div>
  );
};

export default LoginPage;
