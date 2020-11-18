window.BENCHMARK_DATA = {
  "lastUpdate": 1605700769075,
  "repoUrl": "https://github.com/krystal/go-katapult",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "contact@jimeh.me",
            "name": "Jim Myhrberg",
            "username": "jimeh"
          },
          "committer": {
            "email": "contact@jimeh.me",
            "name": "Jim Myhrberg",
            "username": "jimeh"
          },
          "distinct": true,
          "id": "dd681cda878837ba45e1f1756027ad479a61a15e",
          "message": "ci(benchmarks): add benchmark reports to GitHub Actions",
          "timestamp": "2020-11-18T11:44:17Z",
          "tree_id": "0ed4a1ccf27da008754dc3dfe0ce70c12f32ec8c",
          "url": "https://github.com/krystal/go-katapult/commit/dd681cda878837ba45e1f1756027ad479a61a15e"
        },
        "date": 1605700092232,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNameGenerator_RandomHostname",
            "value": 3559,
            "unit": "ns/op",
            "extra": "325344 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "contact@jimeh.me",
            "name": "Jim Myhrberg",
            "username": "jimeh"
          },
          "committer": {
            "email": "contact@jimeh.me",
            "name": "Jim Myhrberg",
            "username": "jimeh"
          },
          "distinct": true,
          "id": "863da59b2c61c24ff3f2b4b5ad6902cabf83fa81",
          "message": "ci(benchmarks): add benchmark reports to GitHub Actions",
          "timestamp": "2020-11-18T11:58:43Z",
          "tree_id": "1b2c69b8c8dac8db948e0e97eec3a71ca9efebaf",
          "url": "https://github.com/krystal/go-katapult/commit/863da59b2c61c24ff3f2b4b5ad6902cabf83fa81"
        },
        "date": 1605700768124,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNameGenerator_RandomHostname",
            "value": 3544,
            "unit": "ns/op",
            "extra": "336682 times\n2 procs"
          },
          {
            "name": "BenchmarkNameGenerator_RandomName_NoPrefix",
            "value": 2569,
            "unit": "ns/op",
            "extra": "467793 times\n2 procs"
          },
          {
            "name": "BenchmarkNameGenerator_RandomName_OnePrefix",
            "value": 2577,
            "unit": "ns/op",
            "extra": "471192 times\n2 procs"
          },
          {
            "name": "BenchmarkNameGenerator_RandomName_TwoPrefixes",
            "value": 2582,
            "unit": "ns/op",
            "extra": "468153 times\n2 procs"
          },
          {
            "name": "BenchmarkNameGenerator_RandomName_ThreePrefixes",
            "value": 2775,
            "unit": "ns/op",
            "extra": "458060 times\n2 procs"
          }
        ]
      }
    ]
  }
}