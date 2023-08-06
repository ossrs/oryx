# SRS Plugin for BT

To create a plugin for [BT](https://bt.cn).

## Usage

To create a `bt-srs_cloud.zip` by:

```bash
bash auto/zip.sh
```

Upload the zip to BT panel, and install SRS.

## Log

Install log saved at `/tmp/srs_cloud_install.log`

```bash
cat /tmp/srs_cloud_install.log
```

When install SRS ready, there should be a file at:

```bash
ls -lh /www/server/panel/plugin/srs_cloud/.bt_ready
```

> Note: If not ready, you're not able to install srs cloud in BT or aaPanel.

