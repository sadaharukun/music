//
//  MusicManager.h
//  MyMusic
//
//  Created by sjjwind on 5/15/15.
//  Copyright (c) 2015 sjjwind. All rights reserved.
//

#import <Foundation/Foundation.h>

typedef enum : NSInteger {
    MMMusicChannelPersonal = 0x0,     // 私人频道
    MMMusicChannelLove,               // 我喜欢
    MMMusicChannelHot,                // 热歌
    MMMusicChannelKTV,                // KTV 歌曲
    MMMusicChannelFamous,             // 成名曲
    MMMusicChannelRamdom,             // 随便听听
    MMMusicChannelNetwork,            // 网络歌曲
    MMMusicChannelTV,                 // 影视歌曲
    MMMusicChannelChinaVioce,         // 中国好声音
    MMMusicChannelClassic,            // 经典老歌
    MMMusicChannel70Age,              // 70后歌曲
    MMMusicChannel80Age,              // 80后歌曲
    MMMusicChannel90Age,              // 90后歌曲
    MMMusicChannelChildren,           // 儿童歌曲
    MMMusicChannelNew,                // 新歌
    MMMusicChannelPopular,            // 流行
    MMMusicChannelLightMusic,         // 轻音乐
    MMMusicChannelFresh,              // 小清新
    MMMusicChannelChineseWind,        // 中国风
    MMMusicChannelRock,               // 摇滚
    MMMusicChannelVideo,              // 电影
    MMMusicChannelFolk,               // 民谣
    MMMusicChannelChinese,            // 华语
    MMMusicChannelEurope,             // 欧美
    MMMusicChannelJanpan,             // 小鬼子
    MMMusicChannelKorea,              // 韩语
    MMMusicChannelCantonese,          // 粤语
    MMMusicChannelHappy,              // 欢快
    MMMusicChannelSlow,               // 舒缓
    MMMusicChannelSad,                // 伤感
    MMMusicChannelReleax,             // 轻松
    MMMusicChannelAlone,              // 寂寞
} MMMusicChannel;

@interface MusicManager : NSObject

+ (instancetype) sharedManager;

- (void)fetchRandomListWithChannel:(MMMusicChannel) channel 
                          complete:(void (^)(int errorCode, NSArray *musicList))completion;

- (void)fetchLoveMusicListWithCompletion:(void (^)(int errorCode, NSArray *musicList))completion;

- (void)fetchListenedMusicListWithCompletion:(void (^)(int errorCode, NSArray *musicList))completion;

- (void)loginByUserName:(NSString *)userName 
               password:(NSString *)password 
             completion:(void (^)(int errorCode))completion;

- (void)searchMusic:(NSString *)keyword 
         completion:(void (^)(int errorCode, NSArray *musicList))completion;

- (void)downloadMusic:(NSInteger)musicId 
             complete:(void (^)(int errorCode, NSString *path))completion;

- (void)downloadCoverImage:(NSInteger)musicId 
                  complete:(void (^)(int errorCode, NSString *path))completion;

- (void)downloadLyric:(NSInteger)musicId 
             complete:(void (^)(int errorCode, NSString *path))completion;

- (void)loveMusic:(NSInteger)musicId 
       loveDegree:(NSInteger)degree;

@end
