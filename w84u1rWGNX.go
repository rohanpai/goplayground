package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bmizerany/pat"
	"github.com/carlhoerberg/go-geoip"
	"github.com/garyburd/redigo/redis"
)
