import React, { useEffect, useState } from "react";

export const Home = (props: { name: string }) => {
  return (
    <div>
      {props.name ? "Hello " + props.name : "User are not authentication"}
    </div>
  );
};
