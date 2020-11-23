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
                window.location.href = new URL(
                  `/oidc/login/microsoft?redirectTo=${encodeURIComponent(
                    window.location.href + "login/done"
                  )}`,
                  process.env.REACT_APP_API_ENDPOINT
                ).toString();
              }}
            >
              Inloggen met Microsoft
            </Button>
            <Button
              raised
              style={{
                backgroundColor: "black",
                color: "white",
                marginTop: 32,
              }}
              className={"LoginBtn"}
              icon={
                <svg
                  viewBox="0 0 533.5 544.3"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M533.5 278.4c0-18.5-1.5-37.1-4.7-55.3H272.1v104.8h147c-6.1 33.8-25.7 63.7-54.4 82.7v68h87.7c51.5-47.4 81.1-117.4 81.1-200.2z"
                    fill="#4285f4"
                  />
                  <path
                    d="M272.1 544.3c73.4 0 135.3-24.1 180.4-65.7l-87.7-68c-24.4 16.6-55.9 26-92.6 26-71 0-131.2-47.9-152.8-112.3H28.9v70.1c46.2 91.9 140.3 149.9 243.2 149.9z"
                    fill="#34a853"
                  />
                  <path
                    d="M119.3 324.3c-11.4-33.8-11.4-70.4 0-104.2V150H28.9c-38.6 76.9-38.6 167.5 0 244.4l90.4-70.1z"
                    fill="#fbbc04"
                  />
                  <path
                    d="M272.1 107.7c38.8-.6 76.3 14 104.4 40.8l77.7-77.7C405 24.6 339.7-.8 272.1 0 169.2 0 75.1 58 28.9 150l90.4 70.1c21.5-64.5 81.8-112.4 152.8-112.4z"
                    fill="#ea4335"
                  />
                </svg>
              }
              onMouseOver={() => setIsHovering(true)}
              onMouseOut={() => setIsHovering(false)}
              onClick={() => {
                window.location.href = new URL(
                  `oidc/login/google?redirectTo=${encodeURIComponent(
                    window.location.href + "login/done"
                  )}`,
                  process.env.REACT_APP_API_ENDPOINT
                ).toString();
              }}
            >
              Inloggen met Google
            </Button>
          </div>

          <span>&copy; {new Date().getFullYear()} de auteurs van Timeterm</span>
        </Elevation>
      </Theme>
    </div>
  );
};

export default LoginPage;
