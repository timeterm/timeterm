import React from "react";
import { Typography } from "@rmwc/typography";
import "@rmwc/typography/styles";
import { Button } from "@rmwc/button";
import { Theme } from "@rmwc/theme";
import { Elevation } from "@rmwc/elevation";
import Logo from "./logo-white.svg";

const LoginPage: React.FC = (props) => {
  return (
    <div
      style={{
        width: "100%",
        height: "100%",
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <Theme use={["primaryBg", "onPrimary"]} wrap>
        <Elevation
          z={16}
          style={{
            borderRadius: 16,
            padding: 32,
            height: "30em",
            width: "17em",
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
          </div>

          <Button
            raised
            style={{ backgroundColor: "white", color: "black" }}
            icon={
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 23 23">
                <path fill="#f35325" d="M1 1h10v10H1z" />
                <path fill="#81bc06" d="M12 1h10v10H12z" />
                <path fill="#05a6f0" d="M1 12h10v10H1z" />
                <path fill="#ffba08" d="M12 12h10v10H12z" />
              </svg>
            }
          >
            Inloggen met Microsoft
          </Button>
        </Elevation>
      </Theme>
    </div>
  );
};

export default LoginPage;
