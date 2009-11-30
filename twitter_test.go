//
// Copyright 2009 Bill Casarin <billcasarin@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package twitter

import "testing"
import "os"
import "fmt"

// dont change this
const kId = 5641609144;

var cache *CacheBackend;

func TestValidStatus(t *testing.T) {
  api := NewApi();
  cache = api.cacheBackend;
  errors := api.GetErrorChannel();

  fmt.Printf("<-api.GetStatus() ...\n");
  status := (<-api.GetUser("jb55")).GetStatus();

  verifyValidStatus(status, t);
  verifyValidUser(status.GetUser(), t);
  getAllApiErrors(errors, t);
}

func TestValidUser(t *testing.T) {
  api := NewApi();
  api.SetCache(cache);
  errors := api.GetErrorChannel();

  fmt.Printf("<-api.GetUserById() ...\n");

  // this should be cached from our first test
  // kId is a status id from my twitter account,
  // and it should have grabbed my user info
  // when it grabbed the status
  hitBefore := cache.hit;
  user := <-api.GetUserById(9918032);
  hitAfter := cache.hit;

  if hitBefore == hitAfter {
    t.Error("Cache was not hit on user lookup, expected: cache lookup");
  }

  verifyValidUser(user, t);
  verifyValidStatus(user.GetStatus(), t);
  getAllApiErrors(errors, t);
}

func TestCacheStoredNotZero(t *testing.T) {
  api := NewApi();
  api.SetCache(cache);
  errors := api.GetErrorChannel();

  stored := cache.store;
  assertGreaterThanZero(stored, "Cache.Store", t);

  getAllApiErrors(errors, t);
}

func TestCacheGet(t *testing.T) {
  api := NewApi();
  api.SetCache(cache);
  errors := api.GetErrorChannel();
  var last, current Status;
  last = nil;

  for i := 0; i < 5; i++ {
    current = <-api.GetStatus(kId);
    if current != last && i != 0 {
      t.Errorf("Cache not storing Status data");
    }
    last = current;
  }

  assertGreaterThanZero(cache.hit, "Cache.Hit", t);

  getAllApiErrors(errors, t);
}

func TestValidFollowerList(t *testing.T) {
  api := NewApi();
  api.SetCache(cache);
  errors := api.GetErrorChannel();
  fmt.Printf("<-api.GetFollowers() ...\n");
  users := <-api.GetFollowers("jb55", 0);
  length := len(users);

  if length <= 1 {
    t.Errorf("len(GetFollowers()) <= 1, got %d expected > 1", length);
  }

  for _, user := range users {
    verifyValidUser(user, t);
    if user.GetStatus().GetId() != 0 {
      verifyValidStatus(user.GetStatus(), t);
    }
  }

  getAllApiErrors(errors, t);
}

func TestValidFriendsList(t *testing.T) {
  api := NewApi();
  api.SetCache(cache);
  errors := api.GetErrorChannel();
  fmt.Printf("<-api.GetFriends() ...\n");
  users := <-api.GetFriends("jb55", 0);
  length := len(users);

  if length <= 1 {
    t.Errorf("len(GetFriends()) <= 1, got %d expected > 1", length);
  }

  for _, user := range users {
    verifyValidUser(user, t);
    if user.GetStatus().GetId() != 0 {
      verifyValidStatus(user.GetStatus(), t);
    }
  }

  getAllApiErrors(errors, t);
}

func TestValidSearchResults(t *testing.T) {
  api := NewApi();
  api.SetCache(cache);
  errors := api.GetErrorChannel();
  fmt.Printf("<-api.SearchSimple() ...\n");
  results := <-api.SearchSimple("#ff");
  length := len(results);

  if length <= 1 {
    t.Errorf("len(SearchSimple()) <= 1, got %d expected > 1", length);
  }

  for _, result := range results {
    verifyValidSearchResult(result, t);
  }

  getAllApiErrors(errors, t);
}

func TestValidPublicTimeLine(t *testing.T) {
  api := NewApi();
  api.SetCache(cache);
  errors := api.GetErrorChannel();
  fmt.Printf("<-api.GetPublicTimeline() ...\n");
  statuses := <-api.GetPublicTimeline();
  length := len(statuses);

  if length <= 1 {
    t.Errorf("len(GetPublicTimeline()) <= 1, got %d expected > 1", length);
  }

  if StatusEqual(statuses[0], statuses[1]) {
    t.Errorf("GetPublicTimeline()[0] == GetPublicTimeline()[1], expected different");
  }

  t.Logf("Number of Statuses retrieved: %d", length);
  for _, status := range statuses {
    verifyValidStatus(status, t);
  }

  getAllApiErrors(errors, t);
}


// Authfile: .twitterauth
// Format: single line, two words
//    username password
func authFromFile() {
  return;
}

func verifyValidUser(u User, t *testing.T) {
  assertGreaterThanZero(u.GetId(), "GetId", t);
  assertNotEmpty(u.GetScreenName(), "GetScreenName", t);
  assertNotEmpty(u.GetName(), "GetName", t);
  assertNotNil(u.GetStatus(), "GetStatus", t);
  assertNotEmpty(u.GetScreenName(), "GetScreenName", t);
}

func verifyValidStatus(s Status, t *testing.T) {
  assertGreaterThanZero(s.GetId(), "GetId", t);
  assertNotEmpty(s.GetCreatedAt(), "GetCreatedAt", t);
  assertNotEmpty(s.GetText(), "GetText", t);
  assertNotNil(s.GetUser(), "GetUser", t);
}

func verifyValidSearchResult(r SearchResult, t *testing.T) {
  assertGreaterThanZero(r.GetId(), "GetId", t);
  assertNotEmpty(r.GetText(), "GetText", t);
}

func getAllApiErrors(errors <-chan os.Error, t *testing.T) {
  if len(errors) == 0 {
    return;
  }
  t.Log("--- Errors generated by GoTwitter START ----");
  for err := range errors {
    t.Log(err.String());
  }
  t.Error("--- Errors generated by GoTwitter END ----");
}

func IsEmpty(s string) bool {
  return len(s) == 0;
}

func assertGreaterThanZero(i int64, fn string, t *testing.T) {
  if i <= 0 {
    t.Errorf("%s is <= 0, got %d expected > 0", fn, i);
  }
}

func assertNotEmpty(s, fn string, t *testing.T) {
  if IsEmpty(s) {
    t.Errorf("%s is empty, expected not empty", fn);
  }
}

func assertNotNil(i interface{}, fn string, t *testing.T) {
  if i == nil {
    t.Errorf("%s is nil, expected not nil", fn);
  }
}

func StatusEqual(a, b Status) bool {
  return a.GetCreatedAt() == b.GetCreatedAt() &&
         a.GetCreatedAtInSeconds() == b.GetCreatedAtInSeconds() &&
         a.GetFavorited() == b.GetFavorited() &&
         a.GetId() == b.GetId() &&
         a.GetInReplyToScreenName() == b.GetInReplyToScreenName() &&
         a.GetInReplyToStatusId() == b.GetInReplyToStatusId() &&
         a.GetInReplyToUserId() == b.GetInReplyToUserId() &&
         a.GetNow() == b.GetNow();
}