Check Expiring Certs
====================

This is a simple program to warn if SSL certificate on specified host is to
expire in a week.

Install
-------

    go get github.com/ei-grad/check-expiring-certs

Check some hosts
----------------

    check-expiring-certs 8.8.8.8:443 [2a02:6b8:a::a]:443 ya.ru:443
