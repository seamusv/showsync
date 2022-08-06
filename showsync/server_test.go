package showsync

import (
	"reflect"
	"testing"
)

func Test_ParseServerOutput(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []Entry
		wantErr bool
	}{
		{
			name: "Parse Sonarr ",
			data: sonarr_output,
			want: []Entry{
				{
					Title:    "The.Simpsons.S17E03.1080p.BluRay.x264-TAXES",
					Status:   "Downloading",
					Protocol: "usenet",
				},
			},
			wantErr: false,
		},
		{
			name: "Parse Radarr ",
			data: radarr_output,
			want: []Entry{
				{
					Title:    "Minions.The.Rise.of.Gru.2022.1080p.WEB-DL.DDP5.1.Atmos.H.264-CMRG",
					Status:   "completed",
					Protocol: "torrent",
				},
				{
					Title:    "Luck.2022.720p.WEB.h264-KOGi",
					Status:   "completed",
					Protocol: "torrent",
				},
				{
					Title:    "Luck.2022.1080p.ATVP.WEBRip.DD5.1.X.264-EVO",
					Status:   "completed",
					Protocol: "torrent",
				},
				{
					Title:    "Minions.The.Rise.Of.Gru.2022.REPACK.1080p.WEB.h264-RUMOUR",
					Status:   "completed",
					Protocol: "torrent",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseServerOutput(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseServer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var radarr_output = []byte(`{
  "page": 1,
  "pageSize": 10,
  "sortKey": "timeleft",
  "sortDirection": "ascending",
  "totalRecords": 4,
  "records": [
    {
      "movieId": 3461,
      "languages": [
        {
          "id": 1,
          "name": "English"
        }
      ],
      "quality": {
        "quality": {
          "id": 3,
          "name": "WEBDL-1080p",
          "source": "webdl",
          "resolution": 1080,
          "modifier": "none"
        },
        "revision": {
          "version": 1,
          "real": 0,
          "isRepack": false
        }
      },
      "customFormats": [],
      "size": 4805530238,
      "title": "Minions.The.Rise.of.Gru.2022.1080p.WEB-DL.DDP5.1.Atmos.H.264-CMRG",
      "sizeleft": 0,
      "status": "completed",
      "trackedDownloadStatus": "warning",
      "trackedDownloadState": "importPending",
      "statusMessages": [
        {
          "title": "Minions.The.Rise.of.Gru.2022.1080p.WEB-DL.DDP5.1.Atmos.H.264-CMRG",
          "messages": [
            "No files found are eligible for import in /downloads/completed/Minions.The.Rise.of.Gru.2022.1080p.WEB-DL.DDP5.1.Atmos.H.264-CMRG"
          ]
        }
      ],
      "downloadId": "D75CDF65C0BC89FBD9282CD0C6D4CD694E51350F",
      "protocol": "torrent",
      "downloadClient": "NL",
      "indexer": "Jackett",
      "outputPath": "/downloads/completed/Minions.The.Rise.of.Gru.2022.1080p.WEB-DL.DDP5.1.Atmos.H.264-CMRG",
      "id": 2048631772
    },
    {
      "movieId": 3473,
      "languages": [
        {
          "id": 1,
          "name": "English"
        }
      ],
      "quality": {
        "quality": {
          "id": 5,
          "name": "WEBDL-720p",
          "source": "webdl",
          "resolution": 720,
          "modifier": "none"
        },
        "revision": {
          "version": 1,
          "real": 0,
          "isRepack": false
        }
      },
      "customFormats": [],
      "size": 3626082855,
      "title": "Luck.2022.720p.WEB.h264-KOGi",
      "sizeleft": 0,
      "status": "completed",
      "trackedDownloadStatus": "warning",
      "trackedDownloadState": "importPending",
      "statusMessages": [
        {
          "title": "Luck.2022.720p.WEB.h264-KOGi",
          "messages": [
            "No files found are eligible for import in /downloads/completed/Luck.2022.720p.WEB.h264-KOGi"
          ]
        }
      ],
      "downloadId": "14BA4D81B0DFBBC239896FB99B878B8928379DD1",
      "protocol": "torrent",
      "downloadClient": "NL",
      "indexer": "Jackett",
      "outputPath": "/downloads/completed/Luck.2022.720p.WEB.h264-KOGi",
      "id": 1199250455
    },
    {
      "movieId": 3473,
      "languages": [
        {
          "id": 1,
          "name": "English"
        }
      ],
      "quality": {
        "quality": {
          "id": 15,
          "name": "WEBRip-1080p",
          "source": "webrip",
          "resolution": 1080,
          "modifier": "none"
        },
        "revision": {
          "version": 1,
          "real": 0,
          "isRepack": false
        }
      },
      "customFormats": [],
      "size": 2044120574,
      "title": "Luck.2022.1080p.ATVP.WEBRip.DD5.1.X.264-EVO",
      "sizeleft": 0,
      "status": "completed",
      "trackedDownloadStatus": "warning",
      "trackedDownloadState": "importPending",
      "statusMessages": [
        {
          "title": "Luck.2022.1080p.ATVP.WEBRip.DD5.1.X.264-EVO",
          "messages": [
            "No files found are eligible for import in /downloads/completed/Luck.2022.1080p.ATVP.WEBRip.DD5.1.X.264-EVO"
          ]
        }
      ],
      "downloadId": "7E08E7D17A84687A380DF055284F5EAD916DBF61",
      "protocol": "torrent",
      "downloadClient": "NL",
      "indexer": "Jackett",
      "outputPath": "/downloads/completed/Luck.2022.1080p.ATVP.WEBRip.DD5.1.X.264-EVO",
      "id": 153899728
    },
    {
      "movieId": 3461,
      "languages": [
        {
          "id": 1,
          "name": "English"
        }
      ],
      "quality": {
        "quality": {
          "id": 3,
          "name": "WEBDL-1080p",
          "source": "webdl",
          "resolution": 1080,
          "modifier": "none"
        },
        "revision": {
          "version": 2,
          "real": 0,
          "isRepack": true
        }
      },
      "customFormats": [],
      "size": 4789677939,
      "title": "Minions.The.Rise.Of.Gru.2022.REPACK.1080p.WEB.h264-RUMOUR",
      "sizeleft": 0,
      "status": "completed",
      "trackedDownloadStatus": "warning",
      "trackedDownloadState": "importPending",
      "statusMessages": [
        {
          "title": "Minions.The.Rise.Of.Gru.2022.REPACK.1080p.WEB.h264-RUMOUR",
          "messages": [
            "No files found are eligible for import in /downloads/completed/Minions.The.Rise.Of.Gru.2022.REPACK.1080p.WEB.h264-RUMOUR"
          ]
        }
      ],
      "downloadId": "01DCF1F9CECF003CE5879933D8E704321449DF33",
      "protocol": "torrent",
      "downloadClient": "NL",
      "indexer": "Jackett",
      "outputPath": "/downloads/completed/Minions.The.Rise.Of.Gru.2022.REPACK.1080p.WEB.h264-RUMOUR",
      "id": 762061305
    }
  ]
}`)

var sonarr_output = []byte(`[
  {
    "series": {
      "title": "The Simpsons",
      "sortTitle": "simpsons",
      "seasonCount": 34,
      "status": "continuing",
      "overview": "Set in Springfield, the average American town, the show focuses on the antics and everyday adventures of the Simpson family; Homer, Marge, Bart, Lisa and Maggie, as well as a virtual cast of thousands. Since the beginning, the series has been a pop culture icon, attracting hundreds of celebrities to guest star. The show has also made name for itself in its fearless satirical take on politics, media and American life in general.",
      "network": "FOX",
      "airTime": "20:00",
      "images": [
        {
          "coverType": "banner",
          "url": "https://artworks.thetvdb.com/banners/graphical/71663-g11.jpg"
        },
        {
          "coverType": "poster",
          "url": "https://artworks.thetvdb.com/banners/posters/71663-15.jpg"
        },
        {
          "coverType": "fanart",
          "url": "https://artworks.thetvdb.com/banners/fanart/original/71663-3.jpg"
        }
      ],
      "seasons": [
        {
          "seasonNumber": 0,
          "monitored": false
        },
        {
          "seasonNumber": 1,
          "monitored": false
        },
        {
          "seasonNumber": 2,
          "monitored": false
        },
        {
          "seasonNumber": 3,
          "monitored": false
        },
        {
          "seasonNumber": 4,
          "monitored": false
        },
        {
          "seasonNumber": 5,
          "monitored": false
        },
        {
          "seasonNumber": 6,
          "monitored": false
        },
        {
          "seasonNumber": 7,
          "monitored": false
        },
        {
          "seasonNumber": 8,
          "monitored": false
        },
        {
          "seasonNumber": 9,
          "monitored": false
        },
        {
          "seasonNumber": 10,
          "monitored": false
        },
        {
          "seasonNumber": 11,
          "monitored": false
        },
        {
          "seasonNumber": 12,
          "monitored": false
        },
        {
          "seasonNumber": 13,
          "monitored": false
        },
        {
          "seasonNumber": 14,
          "monitored": true
        },
        {
          "seasonNumber": 15,
          "monitored": true
        },
        {
          "seasonNumber": 16,
          "monitored": true
        },
        {
          "seasonNumber": 17,
          "monitored": true
        },
        {
          "seasonNumber": 18,
          "monitored": true
        },
        {
          "seasonNumber": 19,
          "monitored": false
        },
        {
          "seasonNumber": 20,
          "monitored": false
        },
        {
          "seasonNumber": 21,
          "monitored": false
        },
        {
          "seasonNumber": 22,
          "monitored": false
        },
        {
          "seasonNumber": 23,
          "monitored": false
        },
        {
          "seasonNumber": 24,
          "monitored": false
        },
        {
          "seasonNumber": 25,
          "monitored": false
        },
        {
          "seasonNumber": 26,
          "monitored": false
        },
        {
          "seasonNumber": 27,
          "monitored": false
        },
        {
          "seasonNumber": 28,
          "monitored": false
        },
        {
          "seasonNumber": 29,
          "monitored": false
        },
        {
          "seasonNumber": 30,
          "monitored": false
        },
        {
          "seasonNumber": 31,
          "monitored": false
        },
        {
          "seasonNumber": 32,
          "monitored": false
        },
        {
          "seasonNumber": 33,
          "monitored": true
        },
        {
          "seasonNumber": 34,
          "monitored": true
        }
      ],
      "year": 1989,
      "path": "/media/Episodes/The Simpsons",
      "profileId": 6,
      "languageProfileId": 1,
      "seasonFolder": true,
      "monitored": true,
      "useSceneNumbering": false,
      "runtime": 25,
      "tvdbId": 71663,
      "tvRageId": 6190,
      "tvMazeId": 83,
      "firstAired": "1989-12-17T08:00:00Z",
      "lastInfoSync": "2022-08-06T03:09:10.161333Z",
      "seriesType": "standard",
      "cleanTitle": "thesimpsons",
      "imdbId": "tt0096697",
      "titleSlug": "the-simpsons",
      "certification": "TV-PG",
      "genres": [
        "Animation",
        "Comedy"
      ],
      "tags": [],
      "added": "2020-03-24T02:38:27.436691Z",
      "ratings": {
        "votes": 24136,
        "value": 8.9
      },
      "qualityProfileId": 6,
      "id": 216
    },
    "episode": {
      "seriesId": 216,
      "episodeFileId": 0,
      "seasonNumber": 17,
      "episodeNumber": 3,
      "title": "Milhouse of Sand and Fog",
      "airDate": "2005-09-25",
      "airDateUtc": "2005-09-26T00:00:00Z",
      "overview": "When Maggie is showing signs of being ill, the family goes to “the more boisterous house of worship” in town to find Dr. Hibbert, who tells them that Maggie is developing the chicken pox.  After Maggie develops the disease, Marge tries to keep Homer away from her, since he has never had them.  After Flanders expresses an interest in getting his kids infected, Homer and Marge open up the house for a “pox party.”  Milhouse’s divorced parents are both at the party and after some “Margerita’s” are consumed, find themselves getting back together.  Meanwhile, Homer has developed the chicken pox and Marge tries to keep him from scratching.  Milhouse likes the idea of his parents getting back together, but then begins to hate it when he has trouble getting either of them to pay any attention to him.  After seeing an episode of The O.C. Milhouse and Bart come up with a plan to get his parent’s separated again, they plant one of Marge’s bras in Kirk & Luann’s bed.  They don’t succeed in breaking them up; rather they break up Homer and Marge.  Even after Bart confesses his guilt, Marge doesn’t want anything to do with Homer, since he obviously doesn’t trust him anymore.  Bart concocts an outrageous scheme to get them back together, but it goes terribly wrong and both he and Homer find themselves in the river heading toward the falls needing to place their trust in Marge.",
      "hasFile": false,
      "monitored": true,
      "absoluteEpisodeNumber": 359,
      "unverifiedSceneNumbering": false,
      "lastSearchTime": "2022-08-06T13:49:40.395066Z",
      "id": 15118
    },
    "quality": {
      "quality": {
        "id": 7,
        "name": "Bluray-1080p",
        "source": "bluray",
        "resolution": 1080
      },
      "revision": {
        "version": 1,
        "real": 0,
        "isRepack": false
      }
    },
    "size": 1663694951.0,
    "title": "The.Simpsons.S17E03.1080p.BluRay.x264-TAXES",
    "sizeleft": 1249198797.0,
    "status": "Downloading",
    "trackedDownloadStatus": "Ok",
    "statusMessages": [],
    "downloadId": "1f051bd9c37a45da8e990e2509cb3a11",
    "protocol": "usenet",
    "id": 726935495
  }
]`)
