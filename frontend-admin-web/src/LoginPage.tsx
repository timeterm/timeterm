import React from "react";
import { Card } from "@rmwc/card";
import "@rmwc/card/styles";
import { Typography } from "@rmwc/typography";
import "@rmwc/typography/styles";
import { Button } from "@rmwc/button";

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
      <Card
        style={{
          padding: 16,
          borderRadius: 8,
          height: "30em",
          justifyContent: "space-between",
          display: "flex",
          flexDirection: "column",
        }}
      >
        <div style={{ display: "flex", flexDirection: "column" }}>
          <Typography use="headline5">Inloggen</Typography>
          Welkom terug!
        </div>

        <Button raised>Inloggen met Microsoft</Button>
      </Card>
    </div>
  );
};

export default LoginPage;
