package track

import "errors"

var ErrNotFound = errors.New("not found")

//SELECT track_1."ID" AS "track_1_ID", author.name AS author_name, author.photo_url AS author_photo_url, author.created_at AS author_created_at, author.updated_at AS author_updated_at, author."ID" AS "author_ID"
//FROM track AS track_1 JOIN track_artist AS track_artist_1 ON track_1."ID" = track_artist_1."track_ID" JOIN author ON author."ID" = track_artist_1."author_ID"
//WHERE track_1."ID" IN ($1::INTEGER)
