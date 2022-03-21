import React from "react";
import {TutorialsToast, useTutorials} from "../components/TutorialsButton";

export default function ScenarioTutorials({urls}) {
  const movieTutorials = useTutorials(React.useRef([
    {author: 'SRS', id: 'BV1844y1L7dL'},
    {author: '徐光磊', id: 'BV1RS4y1G7tb'},
    {author: 'SRS', id: 'BV1Nb4y1t7ij'},
    {author: '瓦全', id: 'BV1SF411t7Li'},
    {author: '王大江', id: 'BV16r4y1q7ZT'},
    {author: '骆合祥', id: 'BV16T4y1U7CN'},
    {author: '周亮', id: 'BV1gT4y1U76d'},
    {author: '崔国栋', id: 'BV1aS4y1G7iG'},
    {author: 'SRS', id: 'BV1KY411V7uc'},
    {author: '唐为', id: 'BV14S4y1k7gr'},
  ]));

  return (
      <TutorialsToast tutorials={movieTutorials} />
  );
}

