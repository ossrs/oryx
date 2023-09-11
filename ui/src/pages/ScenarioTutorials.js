//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {TutorialsToast, useTutorials} from "../components/TutorialsButton";

export default function ScenarioTutorials() {
  const movieTutorials = useTutorials({
    bilibili: React.useRef([
      {author: 'SRS', id: 'BV1844y1L7dL'},
      {author: '徐光磊', id: 'BV1RS4y1G7tb'},
      {author: '程晓龙', id: 'BV1tZ4y1R7qp'},
      {author: 'SRS', id: 'BV1Nb4y1t7ij'},
      {author: '王大江', id: 'BV16r4y1q7ZT'},
      {author: '马景瑞', id: 'BV1c341177e7'},
      {author: '周亮', id: 'BV1gT4y1U76d'},
      {author: '崔国栋', id: 'BV1aS4y1G7iG'},
      {author: 'SRS', id: 'BV1KY411V7uc'},
      {author: '唐为', id: 'BV14S4y1k7gr'},
      {author: 'SRS', id: "BV1cq4y1e7Au"},
      {author: '宝哥', id: 'BV1G3411d7Gb'},
    ]),
    medium: React.useRef([
      {id: 'e9fe6f314ac6'},
      {id: '9748ae754c8c'},
      {id: 'ec18dfae7d6f'},
      {id: 'cb618777639f'},
      {id: '38be22beec57'},
      {id: '2aa792c35b25'},
      {id: 'ba1895828b4f'},
    ])
  });

  return (
      <TutorialsToast tutorials={movieTutorials} />
  );
}

