import React, { useState } from 'react';
import { css } from '@emotion/css';
import { useStyles2, useTheme2, Button } from '@grafana/ui';
import type { GrafanaTheme2 } from '@grafana/data';

interface EditorRowProps {
  label: string;
  collapsible?: boolean;
  title?: () => string | React.ReactNode;
  collapsed?: boolean;
  dataTestId?: string;
  children: React.ReactNode;
}

export const EditorRow: React.FC<EditorRowProps> = ({ label, collapsible, collapsed = true, title, children }) => {
  const styles = useStyles2(getStyles);
  const theme = useTheme2();
  const [show, setShow] = useState(collapsed);
  return (
    <div className={styles.root}>
      {collapsible && (
        <div style={{ display: 'flex' }}>
          <Button
            icon={show ? 'angle-down' : 'angle-right'}
            fill="text"
            size="sm"
            variant="secondary"
            onClick={(e) => {
              setShow(!show);
              e.preventDefault();
            }}
            style={{ marginRight: '10px' }}
          />
          <span
            onClick={(e) => {
              setShow(!show);
              e.preventDefault();
            }}
          >
          <b className={styles.collapseTile}>{label}</b>
          </span>
          <span className={styles.collapseTileSecondary}>{title ? title() : 'Options'}</span>
        </div>
      )}
      {show && (
        <div
          style={{
            display: 'flex',
            flexWrap: 'wrap',
            justifyContent: 'flex-start',
            alignContent: 'flex-start',
            alignItems: 'flex-start',
            gap: theme.spacing(1),
            marginTop: '0px',
            marginLeft: '0px',
            flexDirection: 'row',
            width: '100%',
        }}>
          {children}
        </div>
      )}
    </div>
  );
};

const getStyles = (theme: GrafanaTheme2, width?: number | string, borderColor = 'transparent', horizontal = false) => ({
    root: css({
      minWidth: theme.spacing(width ?? 0),
      // boxShadow: `0px 0px 4px 0px ${theme.colors.border.weak}`,
      border: `1px solid ${theme.colors.border.medium}`,
      padding: theme.spacing('8px'),
      background: theme.colors.background.primary,
      // marginRight: horizontal ? '10px' : '5px',
    }),
    collapseTile: css({
      marginRight: theme.spacing(1),
      color: theme.colors.secondary.text,
    }),
    collapseTileSecondary: css({
      color: theme.colors.text.secondary,
      fontSize: theme.typography.bodySmall.fontSize,
      '&:hover': {
        color: theme.colors.secondary.text,
      },
    }),
});

