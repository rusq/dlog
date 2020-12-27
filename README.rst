====================
 Standard Logger Ex
====================

.. image:: https://travis-ci.com/rusq/dlog.svg?branch=master
    :target: https://travis-ci.com/rusq/dlog

 What is this?
===============

This is a simple wrapper around the standard go runtime logger with only one
goal:

* ADD DEBUG OUTPUT FUNCTIONS.

Functions that are available for the caller:

* All those in "log" AND
* Debug
* Debugf
* Debugln

On the package base level these functions will print output only if the
``DEBUG`` environment variable is present and have some non-empty value.

Otherwise, one can construct a new logger::

  func New(out io.Writer, prefix string, flag int, debug bool) *Logger

for the flags, one could use the values from the standard logger.


Usage
=====

Import it in your project replacing "log" import entry:

.. code:: go

  import (
    log "github.com/rusq/dlog"
  )

Use:

.. code:: go

   log.Debug("hello debug!")
