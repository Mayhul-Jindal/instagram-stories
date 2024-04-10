-- name: FollowUser :exec
INSERT INTO followers (
  follower_id, following_id
) VALUES (
  $1, $2
);

-- name: GetFollowingEmails :many
with following_cte as (
	select following_id from followers
	where follower_id = $1
)
SELECT email from users
join following_cte on users.id = following_cte.following_id
limit $2
offset $3;

-- name: GetFollowersEmails :many
with followers_cte as (
    select follower_id from followers
    where following_id = $1
)
SELECT email from users
join followers_cte on users.id = followers_cte.follower_id
limit $2
offset $3;




-- name: GetFollowingIDs :many
select following_id from followers
where follower_id = $1;

-- name: GetFollowersIDs :many
select follower_id from followers
where following_id = $1;
